import { expect, test } from '@grafana/plugin-e2e';

test.describe('Query Editor', () => {


  test('data query is successful when query is valid', async ({ page, panelEditPage, selectors }) => {
    await panelEditPage.datasource.set('SurrealDB');
    await panelEditPage.setVisualization('Table');

    await page.waitForFunction(() => (window as any).monaco);
    await panelEditPage.getByTestIdOrAriaLabel(selectors.components.CodeEditor.container).click();
    await page.keyboard.press('Meta+A');
    await page.keyboard.press('Control+A');
    await page.keyboard.insertText('SELECT day, sum_sales FROM daily_sales');

    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.getErrorIcon()).not.toBeVisible();
    await expect(panelEditPage.panel.fieldNames).toContainText(['day', 'sum_sales']);
    await expect(panelEditPage.panel.data).toContainText([
      '2023-11-28',
      '558375',
      '2023-07-30',
      '600534',
      '2023-10-13',
      '509895',
      '2023-01-11',
      '547645',
      '2023-11-26',
      '840532',
    ]);
  });

  test('data query fails when query is invalid', async ({ page, panelEditPage, selectors }) => {
    await panelEditPage.datasource.set('SurrealDB');
    await panelEditPage.setVisualization('Table');

    await page.waitForFunction(() => (window as any).monaco);
    await panelEditPage.getByTestIdOrAriaLabel(selectors.components.CodeEditor.container).click();
    await page.keyboard.press('Meta+A');
    await page.keyboard.press('Control+A');
    await page.keyboard.insertText('!SELECT day, sum_sales FROM daily_sales');

    await expect(panelEditPage.refreshPanel()).not.toBeOK();
    await expect(panelEditPage.panel.getErrorIcon()).toBeVisible();
  });
});
