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
    name: 1,
    criteria: isUserInternal
  },
  {
    name: 2,
    criteria: isUserFreeTier
  },
  {
    name: 3,
    criteria: everyone
  }
];

async function getRelease(request) {
  const stateJSON = await RELEASE_STATE.get("state");
  const state = JSON.parse(stateJSON);

  if (!state.nextRelease) {
    return state.currentRelease;
  }

  const currentStageIndex = stages.findIndex(
    stage => stage.name === state.currentStage
  );

  if (
    stages
      .slice(0, currentStageIndex + 1)
      .some(stage => stage.criteria(request))
  ) {
    return state.nextRelease;
  }

  return state.currentRelease;
}

async function handleRequest(request) {
  // const release = await getRelease(request);
  const originalPath = new URL(request.url).pathname.slice(1);
  const path = originalPath === '' ? 'index.html' : originalPath;
  const body = await APP_DEPLOYS.get(`current/${path}`);

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
