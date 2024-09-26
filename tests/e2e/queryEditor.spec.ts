import { expect, test } from '@grafana/plugin-e2e';

test.describe('Query Editor', () => {
  test('data query is successful when query is valid', async ({ page, panelEditPage, selectors }) => {
    await panelEditPage.datasource.set('SurrealDB');
    await panelEditPage.setVisualization('Table');

    await page.waitForFunction(() => (window as any).monaco);
    await panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
    await page.keyboard.press('Meta+A');
    await page.keyboard.press('Control+A');
    await page.keyboard.insertText('SELECT month, sum_sales FROM monthly_sales ORDER BY sum_sales DESC');

    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.getErrorIcon()).not.toBeVisible();
    await expect(panelEditPage.panel.fieldNames).toContainText(['month', 'sum_sales']);
    await expect(panelEditPage.panel.data).toContainText([
      '2023-03',
      '17305',
      '2024-03',
      '17201',
      '2024-10',
      '17046',
      '2024-04',
      '16591',
      '2024-09',
      '16331',
    ]);
  });

  test('data query fails when query is invalid', async ({ page, panelEditPage, selectors }) => {
    await panelEditPage.datasource.set('SurrealDB');
    await panelEditPage.setVisualization('Table');

    await page.waitForFunction(() => (window as any).monaco);
    await panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
    await page.keyboard.press('Meta+A');
    await page.keyboard.press('Control+A');
    await page.keyboard.insertText('!SELECT month, sum_sales FROM monthly_sales ORDER BY sum_sales DESC');

    await expect(panelEditPage.refreshPanel()).not.toBeOK();
    await expect(panelEditPage.panel.getErrorIcon()).toBeVisible();
  });
});
