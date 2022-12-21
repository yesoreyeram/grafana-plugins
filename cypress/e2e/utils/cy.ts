const timeout = 30 * 1000;

export const get = (
  selector: string,
  options: Record<string, any> = { timeout }
) => {
  return cy.get(selector, options);
};
export const contains = (
  selector: string,
  options: Record<string, any> = { timeout }
) => {
  return cy.contains(selector, options);
};
export const visit = (
  selector: string,
  options: Partial<Cypress.VisitOptions> = { timeout }
) => {
  return cy.visit(selector, options);
};
export const log = (message: string) => {
  cy.log(message);
};
export const wait = (ms: number = 1) => {
  cy.wait(ms);
};
