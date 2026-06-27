export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    const { pathname } = url;

    const repo = env.REPO || 'HsukqiLee/MediaUnlockTest';
    const fallbackVer = env.FALLBACK_VERSION || 'v1.8.5-1770436107';
    const token = env.GITHUB_TOKEN || '';

    if (pathname === '/api/ip-info') {
      const data = {
        ip: request.headers.get('cf-connecting-ip') || 'Unknown',
        country: request.cf?.country || 'Unknown',
        region: request.cf?.region || 'Unknown',
        city: request.cf?.city || 'Unknown',
        timezone: request.cf?.timezone || 'Unknown',
        asn: request.cf?.asn || 'Unknown',
        organization: request.cf?.asOrganization || 'Unknown'
      };

      return new Response(JSON.stringify(data, null, 2), {
        status: 200,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
          'Access-Control-Allow-Origin': '*'
        }
      });
    }

    if (pathname.startsWith('/test/latest/version') || pathname.startsWith('/monitor/latest/version')) {
      const latestVersion = await getLatestVersion(repo, token, fallbackVer);
      return new Response(latestVersion.trim(), { status: 200 });
    }

    const regex = /^\/(test|monitor)\/([^/]+)\/(.*)$/;
    const match = pathname.match(regex);
    if (match) {
      const [, type, version, filename] = match;

      let githubUrl;
      if (version === 'latest') {
        githubUrl = `https://github.com/${repo}/releases/latest/download/${filename}`;
      } else {
        githubUrl = `https://github.com/${repo}/releases/download/${version}/${filename}`;
      }

      return await fetch(githubUrl);
    }

    // Otherwise, serve static assets
    return env.ASSETS.fetch(request);
  }
};

async function getLatestVersion(repo, token, fallbackVer) {
  const headers = {
    'Accept': 'application/vnd.github.v3+json',
    'User-Agent': 'MediaUnlockTest-CF-Worker',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(
      `https://api.github.com/repos/${repo}/releases/latest`,
      { headers }
    );

    if (response.ok) {
      const data = await response.json();
      return data.tag_name;
    }
  } catch (e) {
    console.error('GitHub API error:', e);
  }

  try {
    const fallbackResponse = await fetch(
      `https://raw.githubusercontent.com/${repo}/refs/heads/main/VERSION`
    );
    if (fallbackResponse.ok) {
      return (await fallbackResponse.text()).trim();
    }
  } catch (e) {
    console.error('Fallback fetch error:', e);
  }

  return fallbackVer;
}
