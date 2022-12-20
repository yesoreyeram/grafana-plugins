const lookupTimeout = 30 * 1000;

export const login = () => {
  cy.visit("/login", { timeout: lookupTimeout });
  cy.contains("Welcome to Grafana", { timeout: lookupTimeout });
  cy.get(`[aria-label="${"Username input field"}"]`, { timeout: lookupTimeout })
    .should("be.visible")
    .type("grafana");
  cy.get(`[aria-label="${"Password input field"}"]`, { timeout: lookupTimeout })
    .should("be.visible")
    .type("grafana");
  cy.get(`[aria-label="${"Login button"}"]`, {
    timeout: lookupTimeout,
  }).click();
  cy.get(".login-page", { timeout: lookupTimeout }).should("not.exist");
};
