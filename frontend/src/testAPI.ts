import { getExperts } from "./api/api";

(async () => {
  try {
    const experts = await getExperts();
    console.log("Experts:", experts);
  } catch (error) {
    console.error("API Test Error:", error);
  }
})();
