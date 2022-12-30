import { login } from "./flows/login";
import { addDataSource } from "./flows/addDatasource";
import { uuid } from "./utils/uuid";

describe("yesoreyeram-vercel-datasource", () => {
  it("new datasource instance should work without error", () => {
    login();
    addDataSource(
      "Vercel",
      uuid(),
      "error loading config. invalid/empty api token"
    );
  });
});
