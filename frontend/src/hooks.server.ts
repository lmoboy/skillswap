import { checkAuth } from '$lib/api/auth';
import { redirect, type Handle } from '@sveltejs/kit'
import process from 'process'
const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';
const securePath = ['/course', '/settings', '/swapping'];

export const handle: Handle = async ({ event, resolve }) => {
  const { url, request } = event;

  // 1. Handle WebSocket upgrade requests first
  if (request.headers.get('upgrade')?.toLowerCase() === 'websocket') {
    return new Response(null, {
      status: 426,
      statusText: 'Upgrade Required',
      headers: {
        'X-WebSocket-Backend': BACKEND_URL
      }
    });
  }

  // 2. Proxy backend API and uploads (auth not required)
  if (url.pathname.startsWith('/api') || url.pathname.startsWith('/uploads')) {
    const backendUrl = `${BACKEND_URL}${url.pathname}${url.search}`;
    const headers = new Headers(request.headers);
    headers.delete('host');
    headers.delete('connection');

    try {
      const backendResponse = await fetch(backendUrl, {
        method: request.method,
        headers,
        body:
          request.method !== 'GET' && request.method !== 'HEAD'
            ? await request.arrayBuffer()
            : undefined,
        duplex: 'half'
      } as RequestInit);

      const responseHeaders = new Headers(backendResponse.headers);
      responseHeaders.delete('transfer-encoding');

      // Add CORS headers
      responseHeaders.set(
        'Access-Control-Allow-Origin',
        request.headers.get('origin') || 'localhost:8080'
      );
      responseHeaders.set('Access-Control-Allow-Credentials', 'true');
      responseHeaders.set(
        'Access-Control-Allow-Methods',
        'GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD'
      );
      responseHeaders.set('Access-Control-Allow-Headers', '*');
      responseHeaders.set('Access-Control-Expose-Headers', '*');

      return new Response(backendResponse.body, {
        status: backendResponse.status,
        statusText: backendResponse.statusText,
        headers: responseHeaders
      });
    } catch (error) {
      console.error('Backend proxy error:', error);
      return new Response(JSON.stringify({ error: 'Backend unavailable' }), {
        status: 502,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin':
            request.headers.get('origin') || 'localhost:8080',
          'Access-Control-Allow-Credentials': 'true'
        }
      });
    }
  }

  // 3. Auth check for restricted routes
  try {
    const isUser = await checkAuth();
    const isPublic = !securePath.some(path => url.pathname.startsWith(path));

    // If user not logged in and trying to access a restricted route
    if (!isUser && !isPublic) {
      throw redirect(303, '/auth/login');
    }
  } catch (error) {
    console.error('Auth check failed:', error);
    throw redirect(308, '/auth/login');
  }

  // 4. Continue to normal SvelteKit handling
  const response = await resolve(event);

  // 5. Add CORS headers to all responses
  response.headers.set(
    'Access-Control-Allow-Origin',
    request.headers.get('origin') || 'localhost:8080'
  );
  response.headers.set('Access-Control-Allow-Credentials', 'true');
  response.headers.set(
    'Access-Control-Allow-Methods',
    'GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD'
  );
  response.headers.set('Access-Control-Allow-Headers', '*');
  response.headers.set('Access-Control-Expose-Headers', '*');

  return response;
};