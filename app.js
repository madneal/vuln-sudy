const puppeteer = require('puppeteer');

const url = 'http://127.0.0.1/dvwa/login.php';
const username = 'admin';
const pass = 'password';

(async () => {
  const browser = await puppeteer.launch({
    headless: false
  });
  const page = await browser.newPage();
  await page.goto(url);
  await page.foucus('.loginInput:nth-child(1)');
  await page.type(username);
  await page.focus('.loginInput:nth-child(2)');
  await page.type(pass);
  await page.click('.submit > input');
  await page.waitForNavigation();

  await page.click('#main_menu_padded > ul:nth-child(3) > li');
  await page.focus('.vulnerable_code_area > input:nth-child(1)');
  await page.type('11341234');
  await page.focus('.vulnerable_code_area > input:nth-child(2)');
  await page.type('43213423412342134');
  await page.click('.vulnerable_code_area > input:nth-child(3)');

})();