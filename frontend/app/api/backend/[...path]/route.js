const BACKEND_BASE_URL = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

async function forward(request, paramsPromise) {
  const params = await paramsPromise;
  const path = (params?.path || []).join("/");
  const url = `${BACKEND_BASE_URL}/api/v1/${path}${request.nextUrl.search || ""}`;

  const headers = new Headers();
  const auth = request.headers.get("authorization");
  if (auth) {
    headers.set("authorization", auth);
  }
  headers.set("content-type", request.headers.get("content-type") || "application/json");

  const method = request.method.toUpperCase();
  const hasBody = !["GET", "HEAD"].includes(method);
  const body = hasBody ? await request.arrayBuffer() : undefined;

  let upstream;
  try {
    upstream = await fetch(url, {
      method,
      headers,
      body,
      cache: "no-store"
    });
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown upstream error";
    return new Response(
      JSON.stringify({
        error: "Backend is unavailable",
        backend_url: BACKEND_BASE_URL,
        detail: errorMessage
      }),
      {
        status: 502,
        headers: { "content-type": "application/json; charset=utf-8" }
      }
    );
  }

  const text = await upstream.text();
  return new Response(text, {
    status: upstream.status,
    headers: {
      "content-type": upstream.headers.get("content-type") || "text/plain; charset=utf-8"
    }
  });
}

export async function GET(request, { params }) {
  return forward(request, params);
}

export async function POST(request, { params }) {
  return forward(request, params);
}

export async function PUT(request, { params }) {
  return forward(request, params);
}

export async function DELETE(request, { params }) {
  return forward(request, params);
}
