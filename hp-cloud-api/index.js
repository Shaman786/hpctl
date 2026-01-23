import express from "express";
import shell from "shelljs";
import cors from "cors";

// --- CONFIGURATION ---
const app = express();
const PORT = 5000;

// --- MIDDLEWARE ---
app.use(cors());
app.use(express.json());

// --- HELPER: DATABASE CHECK ---
const getServerStatus = (name) => {
  const check = shell.exec(`docker inspect ${name}`, { silent: true });
  return check.code === 0 ? "ACTIVE" : "NOT_FOUND";
};

// --- ROUTES ---

// 1. HEALTH CHECK
app.get("/status", (req, res) => {
  res.json({ system: "Host-Palace Cloud API", status: "operational", time: new Date() });
});

// 2. PROVISION SERVER (POST)
app.post("/api/v1/servers", (req, res) => {
  const { name, image } = req.body;
  console.log(`[OPS] Provisioning Request: ${name} (${image})`);

  if (!name || !image) {
    return res.status(400).json({ error: 'Payload missing "name" or "image".' });
  }

  if (getServerStatus(name) === "ACTIVE") {
    return res.status(409).json({ error: `Server '${name}' already exists.` });
  }

  // Logic: Nginx gets port 80, others get sleep infinity
  const cmd = image.includes("nginx")
    ? `docker run -d --name ${name} -p 80:80 ${image}`
    : `docker run -d --name ${name} ${image} sh -c "sleep infinity"`;

  if (shell.exec(cmd, { silent: true }).code !== 0) {
    return res.status(500).json({ error: "Hypervisor failed to allocate resources." });
  }

  res.status(201).json({
    message: "Instance provisioned successfully",
    server_id: `hp-${Math.random().toString(36).substr(2, 9)}`,
    status: "running",
  });
});

// 3. GET LOGS (NEW ROUTE - ADDED THIS)
app.get("/api/v1/servers/:name/logs", (req, res) => {
  const { name } = req.params;
  console.log(`[OPS] Fetching Logs: ${name}`);

  // Fetch last 50 lines of logs from Docker
  const output = shell.exec(`docker logs --tail 50 ${name}`, { silent: true });

  // If command failed, server likely doesn't exist
  if (output.code !== 0) {
    return res.status(404).json({ error: `Server '${name}' not found.` });
  }

  // Docker logs usually come in stdout or stderr depending on the app
  res.json({ logs: output.stdout + output.stderr });
});

// 4. LIST INFRASTRUCTURE (GET)
app.get("/api/v1/servers", (req, res) => {
  console.log(`[OPS] Fetching Infrastructure List`);
  const cmd = 'docker ps --format "{{.Names}}|{{.Image}}|{{.Status}}"';
  const output = shell.exec(cmd, { silent: true }).stdout;

  const servers = output.trim().split("\n").filter((line) => line).map((line) => {
    const [name, img, status] = line.split("|");
    return { name, image: img, status };
  });

  res.json({ count: servers.length, servers });
});

// 5. DECOMMISSION SERVER (DELETE)
app.delete("/api/v1/servers/:name", (req, res) => {
  const { name } = req.params;
  console.log(`[OPS] Termination Request: ${name}`);

  if (shell.exec(`docker rm -f ${name}`, { silent: true }).code !== 0) {
    return res.status(404).json({ error: `Server '${name}' not found.` });
  }

  res.json({ message: `Instance '${name}' terminated.` });
});

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