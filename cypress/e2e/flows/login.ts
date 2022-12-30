import { get, contains, visit } from "../utils/cy";

const selector = {
  URL: `/login`,
  TextCheck: `Welcome to Grafana`,
  UserNameInputField: `[aria-label="Username input field"]`,
  PasswordInputField: `[aria-label="Password input field"]`,
  LoginButton: `[aria-label="Login button"]`,
  LoginPage: ".login-page",
};

export const login = (username = "grafana", password = "grafana") => {
  cy.viewport(1792, 1017);
  visit(selector.URL);
  contains(selector.TextCheck);
  get(selector.UserNameInputField).should("be.visible").type(username);
  get(selector.PasswordInputField).should("be.visible").type(password);
  get(selector.LoginButton).click();
  get(selector.LoginPage).should("not.exist");
};
