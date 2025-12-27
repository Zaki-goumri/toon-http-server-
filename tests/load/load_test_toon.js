import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 50 },
    { duration: "1m", target: 50 },
    { duration: "30s", target: 100 },
    { duration: "1m", target: 100 },
    { duration: "30s", target: 0 },
  ],
};

export default function () {
  const baseUrl = "http://localhost:8080";

  let res = http.get(`${baseUrl}/users`);
  check(res, {
    "GET all users - status 200": (r) => r.status === 200,
  });

  const createPayload = "name: TestUser\nemail: test@example.com";
  res = http.post(`${baseUrl}/users`, createPayload, {
    headers: { "Content-Type": "application/toon" },
  });
  check(res, {
    "POST create user - status 201": (r) => r.status === 201,
  });

  const lines = res.body.split("\n");
  let userId = null;
  for (let line of lines) {
    if (line.includes("id:")) {
      userId = line.split(":")[1].trim();
      break;
    }
  }

  if (userId) {
    res = http.get(`${baseUrl}/users/${userId}`);
    check(res, {
      "GET single user - status 200": (r) => r.status === 200,
    });

    const updatePayload = "name: UpdatedUser\nemail: updated@example.com";
    res = http.put(`${baseUrl}/users/${userId}`, updatePayload, {
      headers: { "Content-Type": "application/toon" },
    });
    check(res, {
      "PUT update user - status 200": (r) => r.status === 200,
    });

    res = http.del(`${baseUrl}/users/${userId}`);
    check(res, {
      "DELETE user - status 204": (r) => r.status === 204,
    });
  }

  sleep(1);
}
