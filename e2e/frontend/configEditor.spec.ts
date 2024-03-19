import { expect, test } from '@grafana/plugin-e2e';

const PLUGIN_UID = 'grafana-surrealdb-datasource';
const SURREAL_DB_URL = 'ws://surrealdb:8000/rpc';

test.describe('Config Editor', () => {
  test('invalid credentials should return an error', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: PLUGIN_UID });
    await page.getByPlaceholder('ws://localhost:8000/rpc').fill(SURREAL_DB_URL);
    await expect(configPage.saveAndTest()).not.toBeOK();
  });

  test('valid credentials should return a 200 status code', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: PLUGIN_UID });
    configPage.mockHealthCheckResponse({ status: 200 });

    await page.getByPlaceholder('ws://localhost:8000/rpc').fill(SURREAL_DB_URL);
    await page.getByPlaceholder('Database name').fill('test');
    await page.getByPlaceholder('Namespace').fill('test');
    await page.getByPlaceholder('Username').fill('root');
    await page.getByPlaceholder('Password').fill('test');

    await expect(configPage.saveAndTest()).toBeOK();
  });

  test('valid credentials should display a success alert on the page', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: PLUGIN_UID });

    await page.getByPlaceholder('ws://localhost:8000/rpc').fill(SURREAL_DB_URL);
    await page.getByPlaceholder('Database name').fill('test');
    await page.getByPlaceholder('Namespace').fill('test');
    await page.getByPlaceholder('Username').fill('root');
    await page.getByPlaceholder('Password').fill('test');

    await configPage.saveAndTest();
    await expect(configPage).toHaveAlert('success', { hasNotText: 'Datasource updated' });

    await page.pause();
  });

  test('mandatory fields should show error if left empty', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: PLUGIN_UID });

    await page.getByLabel('Database name').fill('');
    await page.keyboard.press('Tab');
    await expect(page.getByText('Database name is required')).toBeVisible();

    await page.getByLabel('Namespace').fill('');
    await page.keyboard.press('Tab');
    await expect(page.getByText('Namespace is required')).toBeVisible();

    await page.getByLabel('Username').fill('');
    await page.keyboard.press('Tab');
    await expect(page.getByText('Username is required')).toBeVisible();

    await expect(configPage.saveAndTest()).not.toBeOK();
  });
});
