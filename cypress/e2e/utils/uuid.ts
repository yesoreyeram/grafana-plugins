export const uuid = () => {
  const uniqueSeed = Date.now().toString();
  return Cypress._.uniqueId(uniqueSeed);
};
