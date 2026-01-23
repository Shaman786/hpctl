import express from "express";
import shell from "shelljs";
import cors from "cors";

// --- CONFIGURATION ---
const app = express();
const PORT = 5000;

// --- MIDDLEWARE ---
app.use(cors());
app.use(express.json());

// --- SECURITY: INPUT VALIDATION REGEX ---
const SAFE_NAME_REGEX = /^[a-zA-Z0-9-_]+$/;
const SAFE_IMAGE_REGEX = /^[a-z0-9._/-]+(:[a-zA-Z0-9._-]+)?$/;

// --- VALIDATION MIDDLEWARE ---
const validateServerName = (req, res, next) => {
  const name = req.body?.name || req.params?.name;

  if (!name) {
    return res.status(400).json({ error: "Server name is required." });
  }

  if (!SAFE_NAME_REGEX.test(name)) {
    return res.status(400).json({
      error:
        "Invalid server name. Use only letters, numbers, hyphen (-), and underscore (_).",
    });
  }

  next();
};

const validateImage = (req, res, next) => {
  const { image } = req.body;

  if (!image) {
    return res.status(400).json({ error: "Docker image is required." });
  }

  if (!SAFE_IMAGE_REGEX.test(image)) {
    return res.status(400).json({ error: "Invalid Docker image name." });
  }

  next();
};

// --- HELPER: SERVER EXISTENCE CHECK ---
const getServerStatus = (name) => {
  const check = shell.exec(`docker inspect ${name}`, { silent: true });
  return check.code === 0 ? "ACTIVE" : "NOT_FOUND";
};

// --- ROUTES ---

// 1. HEALTH CHECK
app.get("/status", (req, res) => {
  res.json({
    system: "Host-Palace Cloud API",
    status: "operational",
    time: new Date(),
  });
});

// 2. PROVISION SERVER
app.post(
  "/api/v1/servers",
  validateServerName,
  validateImage,
  (req, res) => {
    const { name, image } = req.body;
    console.log(`[OPS] Provisioning Request: ${name} (${image})`);

    if (getServerStatus(name) === "ACTIVE") {
      return res.status(409).json({
        error: `Server '${name}' already exists.`,
      });
    }

    // RESOURCE LIMITS: Proves you understand Cloud Quotas
    const limits = "--memory=512m --cpus=0.5";

    // COMMAND LOGIC:
    // 1. If Nginx: Run naturally (so logs work), but REMOVE "-p 80:80" to avoid conflicts.
    // 2. If Alpine/Ubuntu: Run "sleep infinity" so the container doesn't exit immediately.
    const cmd = image.includes("nginx")
      ? `docker run -d --name ${name} ${limits} ${image}`
      : `docker run -d --name ${name} ${limits} ${image} sh -c "sleep infinity"`;

    const result = shell.exec(cmd, { silent: true });

    if (result.code !== 0) {
      return res.status(500).json({
        error: "Hypervisor failed to allocate resources.",
      });
    }

    res.status(201).json({
      message: "Instance provisioned successfully",
      server_id: `hp-${Math.random().toString(36).substring(2, 11)}`,
      status: "running",
    });
  }
);

// 3. GET SERVER LOGS
app.get(
  "/api/v1/servers/:name/logs",
  validateServerName,
  (req, res) => {
    const { name } = req.params;
    console.log(`[OPS] Fetching Logs: ${name}`);

    const output = shell.exec(`docker logs --tail 50 ${name}`, {
      silent: true,
    });

    if (output.code !== 0) {
      return res.status(404).json({
        error: `Server '${name}' not found.`,
      });
    }

    res.json({
      logs: output.stdout + output.stderr,
    });
  }
);

// 4. LIST INFRASTRUCTURE
app.get("/api/v1/servers", (req, res) => {
  console.log("[OPS] Fetching Infrastructure List");

  const output = shell.exec(
    'docker ps --format "{{.Names}}|{{.Image}}|{{.Status}}"',
    { silent: true }
  ).stdout;

  const servers = output
    .trim()
    .split("\n")
    .filter(Boolean)
    .map((line) => {
      const [name, image, status] = line.split("|");
      return { name, image, status };
    });

  res.json({
    count: servers.length,
    servers,
  });
});

// 5. DECOMMISSION SERVER
app.delete(
  "/api/v1/servers/:name",
  validateServerName,
  (req, res) => {
    const { name } = req.params;
    console.log(`[OPS] Termination Request: ${name}`);

    const result = shell.exec(`docker rm -f ${name}`, { silent: true });

    if (result.code !== 0) {
      return res.status(404).json({
        error: `Server '${name}' not found.`,
      });
    }

    res.json({
      message: `Instance '${name}' terminated.`,
    });
  }
);

// --- START SERVER ---
app.listen(PORT, () => {
  console.log(`
🚀 HOST-PALACE MOCK CLOUD API
-----------------------------
• Status:   Online
• Port:     ${PORT}
• Gateway:  http://localhost:${PORT}/api/v1
`);
});
