const isUserInternal = request => {
  const userId = parseInt(request.headers.get("X-User-Id"));
  return userId < 100;
};

const isUserFreeTier = request => {
  const userId = parseInt(request.headers.get("X-User-Id"));
  return userId >= 100 && userId < 200;
};

const everyone = request => true;

const stages = [
  {
    name: '1',
    criteria: isUserInternal
  },
  {
    name: '2',
    criteria: isUserFreeTier
  },
  {
    name: '3',
    criteria: everyone
  }
];

async function getDeployId(request) {
  const stateJSON = await RELEASE_STATE.get("state");
  const state = JSON.parse(stateJSON);

  if (!state.next) {
    return state.current;
  }

  const currentStageIndex = stages.findIndex(
    stage => stage.name === state.stage
  );

  if (
    stages
      .slice(0, currentStageIndex + 1)
      .some(stage => stage.criteria(request))
  ) {
    return state.next;
  }

  return state.current;
}

async function handleRequest(request) {
  const deployId = await getDeployId(request);
  const originalPath = new URL(request.url).pathname.slice(1);
  const path = originalPath === '' ? 'index.html' : originalPath;
  const body = await APP_DEPLOYS.get(`${deployId}/${path}`);

  const extensionsToContentTypes = {
    'css': 'text/css',
    'html': 'text/html',
    'js': 'application/javascript'
  };

  const contentType = extensionsToContentTypes[path.split('.').pop()];

  return new Response(body, {
    headers: { "Content-Type": contentType }
  });
}

addEventListener("fetch", event => {
  event.respondWith(handleRequest(event.request));
});
