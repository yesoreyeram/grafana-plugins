import { login } from "./flows/login";
import {
  addDataSource,
  validateProvisionedDatasource,
} from "./flows/addDatasource";
import { uuid } from "./utils/uuid";

describe("yesoreyeram-hello-datasource", () => {
  it("new datasource instance should work without error", () => {
    login();
    addDataSource("Hello", uuid(), "invalid token");
  });
  it("provisioned datasources should work without error", () => {
    login();
    validateProvisionedDatasource("Hello", "OK", ["Name Transform Mode"]);
  });
});
