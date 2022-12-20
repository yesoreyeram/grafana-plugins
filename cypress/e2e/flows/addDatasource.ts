const lookupTimeout = 30 * 1000;

export const addDataSource = (
  pluginName: string = "TestData DB",
  name = new Date().getTime().toLocaleString(),
  expectedMessage = "Data source is working",
  timeout = 30 * 1000
) => {
  cy.log(`Adding datasource ${pluginName} - ${name}`);
  cy.visit("/datasources/new", { timeout: lookupTimeout });
  cy.contains("Add data source", { timeout: lookupTimeout });
  cy.get(`[aria-label$="data source ${pluginName}"]`, {
    timeout: lookupTimeout,
  })
    .contains(pluginName)
    .scrollIntoView()
    .should("be.visible")
    .click();
  cy.contains(`Type: ${pluginName}`, { timeout: lookupTimeout });
  cy.get(`[aria-label="Data source settings page name input field"]`, {
    timeout: lookupTimeout,
  })
    .clear()
    .type(`${pluginName}-${name}`);
  cy.get(`[aria-label="Data source settings page Save and Test button"]`, {
    timeout: lookupTimeout,
  }).click();
  cy.get(`[aria-label="Data source settings page Alert"]`, {
    timeout: lookupTimeout,
  })
    .should("exist")
    .contains(expectedMessage, { timeout });
};

export const validateProvisionedDatasource = (
  name = "testdata",
  expectedMessage = "Data source is working",
  checks: string[] = ["Name"],
  timeout = 30 * 1000
) => {
  cy.log(`Validating provisioned datasource - ${name}`);
  cy.visit("/datasources", { timeout: lookupTimeout });
  cy.get("h2 a", { timeout: lookupTimeout }).each(($v) => {
    if ($v.text() === name) {
      cy.wrap($v).contains(name).scrollIntoView().should("be.visible").click();
    }
  });
  cy.wait(2 * 1000);
  (checks || []).forEach((c) => {
    cy.contains(c, { timeout });
  });
  cy.get("button:last-child", { timeout: lookupTimeout })
    .contains("Test")
    .scrollIntoView()
    .should("be.visible")
    .click();
  cy.get(`[aria-label="Data source settings page Alert"]`, { timeout })
    .should("exist")
    .contains(expectedMessage, { timeout });
};
