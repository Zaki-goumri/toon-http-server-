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
  const baseUrl = "http://localhost:8081/json";

  let res = http.get(`${baseUrl}/users`);
  check(res, {
    "GET all users - status 200": (r) => r.status === 200,
  });

  const createPayload = JSON.stringify({
    name: "TestUser",
    email: "test@example.com",
  });
  res = http.post(`${baseUrl}/users`, createPayload, {
    headers: { "Content-Type": "application/json" },
  });
  check(res, {
    "POST create user - status 201": (r) => r.status === 201,
  });

  let userId = null;
  try {
    const userData = JSON.parse(res.body);
    userId = userData.id || userData.ID;
  } catch (e) {}

  if (userId) {
    res = http.get(`${baseUrl}/users/${userId}`);
    check(res, {
      "GET single user - status 200": (r) => r.status === 200,
    });

    const updatePayload = JSON.stringify({
      name: "UpdatedUser",
      email: "updated@example.com",
    });
    res = http.put(`${baseUrl}/users/${userId}`, updatePayload, {
      headers: { "Content-Type": "application/json" },
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
