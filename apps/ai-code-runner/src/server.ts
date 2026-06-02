import { createServer, IncomingMessage, ServerResponse } from "node:http";
import { runFeatureJob } from "./jobs/runFeatureJob.js";

const port = Number(process.env.PORT || 8080);

function sendJson(response: ServerResponse, statusCode: number, body: unknown): void {
  response.writeHead(statusCode, { "Content-Type": "application/json" });
  response.end(JSON.stringify(body));
}

async function readJson(request: IncomingMessage): Promise<unknown> {
  const chunks: Buffer[] = [];
  for await (const chunk of request) {
    chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk));
  }

  const raw = Buffer.concat(chunks).toString("utf8");
  if (!raw.trim()) {
    return {};
  }

  return JSON.parse(raw);
}

const server = createServer(async (request, response) => {
  try {
    if (request.method === "GET" && request.url === "/health") {
      sendJson(response, 200, { ok: true, service: "ai-code-runner" });
      return;
    }

    if (request.method === "POST" && request.url === "/jobs/feature") {
      const body = await readJson(request);
      const result = await runFeatureJob(body);
      sendJson(response, result.success ? 200 : 422, result);
      return;
    }

    sendJson(response, 404, { success: false, error: "Not found" });
  } catch (error) {
    sendJson(response, 500, {
      success: false,
      stage: "server",
      error: error instanceof Error ? error.message : "Unknown server error",
      logs: "",
      recommendation: "Check the request payload and runner logs."
    });
  }
});

server.listen(port, () => {
  console.log(`ai-code-runner listening on ${port}`);
});
