import { login } from "./flows/login";
import {
  addDataSource,
  validateProvisionedDatasource,
} from "./flows/addDatasource";
import { uuid } from "./utils/uuid";

describe("yesoreyeram-blank-datasource", () => {
  it("new datasource instance should work without error", () => {
    login();
    addDataSource(
      "Blank",
      uuid(),
      "blank datasource just works but does nothing"
    );
  });
  it("provisioned datasources should work without error", () => {
    login();
    validateProvisionedDatasource(
      "Blank",
      "blank datasource just works but does nothing",
      ["Blank Config Editor"]
    );
  });
});
