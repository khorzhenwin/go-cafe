const BACKEND_BASE_URL = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export async function GET(request) {
  const url = `${BACKEND_BASE_URL}/api/v1/discovery/cafes/static-map${request.nextUrl.search || ""}`;

  let upstream;
  try {
    upstream = await fetch(url, { cache: "no-store" });
  } catch (error) {
    return new Response(error instanceof Error ? error.message : "Static map upstream failed", {
      status: 502,
      headers: { "content-type": "text/plain; charset=utf-8" }
    });
  }

  const body = await upstream.arrayBuffer();
  return new Response(body, {
    status: upstream.status,
    headers: {
      "content-type": upstream.headers.get("content-type") || "image/png",
      "cache-control": upstream.headers.get("cache-control") || "public, max-age=300"
    }
  });
}
