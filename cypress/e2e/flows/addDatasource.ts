import { get, contains, visit, log, wait } from "../utils/cy";

export const addDataSource = (
  pluginName: string = "TestData DB",
  name = new Date().getTime().toLocaleString(),
  expectedMessage = "Data source is working",
  timeout = 30 * 1000
) => {
  log(`Adding datasource ${pluginName} - ${name}`);
  visit("/datasources/new");
  contains("Add data source");
  get(`[aria-label$="data source ${pluginName}"]`)
    .contains(pluginName)
    .scrollIntoView()
    .should("be.visible")
    .click();
  contains(`Type: ${pluginName}`);
  get(`[aria-label="Data source settings page name input field"]`)
    .clear()
    .type(`${pluginName}-${name}`);
  get(`[aria-label="Data source settings page Save and Test button"]`).click();
  get(`[aria-label="Data source settings page Alert"]`)
    .should("exist")
    .contains(expectedMessage, { timeout });
};

export const validateProvisionedDatasource = (
  name = "testdata",
  expectedMessage = "Data source is working",
  checks: string[] = ["Name"],
  timeout = 30 * 1000
) => {
  log(`Validating provisioned datasource - ${name}`);
  visit("/datasources");
  get("h2 a").each(($v) => {
    if ($v.text() === name) {
      cy.wrap($v).contains(name).scrollIntoView().should("be.visible").click();
    }
  });
  wait(2 * 1000);
  (checks || []).forEach((c) => contains(c));
  get("button:last-child")
    .contains("Test")
    .scrollIntoView()
    .should("be.visible")
    .click();
  get(`[aria-label="Data source settings page Alert"]`)
    .should("exist")
    .contains(expectedMessage, { timeout });
};
