import { login } from "./flows/login";
import {
  addDataSource,
  validateProvisionedDatasource,
} from "./flows/addDatasource";
import { uuid } from "./utils/uuid";

describe("yesoreyeram-petstore-datasource", () => {
  it("new datasource instance should work without error", () => {
    login();
    addDataSource("Pet Store", uuid(), "Pet Store datasource plugin works");
  });
  it("provisioned datasources should work without error", () => {
    login();
    validateProvisionedDatasource(
      "Pet Store",
      "Pet Store datasource plugin works",
      ["Pet Store Config Editor"]
    );
  });
});
