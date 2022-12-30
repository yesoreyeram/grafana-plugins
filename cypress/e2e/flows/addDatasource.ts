import { get, contains, visit, log, wait } from "../utils/cy";

const selector = {
  DSList: {
    Url: "/datasources",
    DSTitle: "h2 a",
  },
  DSEdit: {
    TestButton: "button:last-child",
    ValidationMessage: `[aria-label="Data source settings page Alert"]`,
  },
  DSNew: {
    Url: "/datasources/new",
    PluginTile: (pluginName: string) =>
      `[aria-label$="data source ${pluginName}"]`,
    NameInput: `[aria-label="Data source settings page name input field"]`,
    SaveAndTestButton: `[aria-label="Data source settings page Save and Test button"]`,
    ValidationMessage: `[aria-label="Data source settings page Alert"]`,
  },
};

export const addDataSource = (
  pluginName: string = "TestData DB",
  name = new Date().getTime().toLocaleString(),
  expectedMessage = "Data source is working",
  timeout = 30 * 1000
) => {
  log(`Adding datasource ${pluginName} - ${name}`);
  visit(selector.DSNew.Url);
  contains("Add data source");
  get(selector.DSNew.PluginTile(pluginName))
    .contains(pluginName)
    .scrollIntoView()
    .should("be.visible")
    .click();
  contains(`Type: ${pluginName}`);
  get(selector.DSNew.NameInput).clear().type(`${pluginName}-${name}`);
  get(selector.DSNew.SaveAndTestButton).click();
  get(selector.DSNew.ValidationMessage)
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
  visit(selector.DSList.Url);
  get(selector.DSList.DSTitle).each(($v) => {
    if ($v.text() === name) {
      cy.wrap($v).contains(name).scrollIntoView().should("be.visible").click();
    }
  });
  wait(2 * 1000);
  (checks || []).forEach((c) => contains(c));
  get(selector.DSEdit.TestButton)
    .contains("Test")
    .scrollIntoView()
    .should("be.visible")
    .click();
  get(selector.DSEdit.ValidationMessage)
    .should("exist")
    .contains(expectedMessage, { timeout });
};
