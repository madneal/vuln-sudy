const puppeteer = require('puppeteer');

const url = 'http://127.0.0.1/dvwa/login.php';
const username = 'admin';
const pass = 'password';

(async () => {
  const browser = await puppeteer.launch({
    headless: false,
    args: ['--no-sandbox']
  });
  const page = await browser.newPage();
  await page.goto(url);
  await page.click('#content > form > fieldset > input:nth-child(2)');
//  await page.foucus('#content > form > fieldset > input:nth-child(2)');
  await page.type('#content > form > fieldset > input:nth-child(2)','admin');
  await page.click('#content > form > fieldset > input:nth-child(5)');
//  await page.focus('ontent > form > fieldset > input:nth-child(3)');
  await page.type('#content > form > fieldset > input:nth-child(5)', pass);
  await page.click('.submit > input');
  await page.waitForNavigation();

  await page.click('#main_menu_padded > ul:nth-child(2) > li');
  await page.click('#main_body > div > div > form > input[type="text"]:nth-child(2)');
  await page.type('#main_body > div > div > form > input[type="text"]:nth-child(2)', '11341234');
  await page.click('#main_body > div > div > form > input[type="password"]:nth-child(5)');
  await page.type('#main_body > div > div > form > input[type="password"]:nth-child(5)', '43213423412342134');
  await page.click('#main_body > div > div > form > input[type="submit"]:nth-child(8)');
  await page.click('#main_body > div > div > form > input[type="text"]:nth-child(2)');
  await page.type('#main_body > div > div > form > input[type="text"]:nth-child(2)', '11341234');
  await page.click('#main_body > div > div > form > input[type="password"]:nth-child(5)');
  await page.type('#main_body > div > div > form > input[type="password"]:nth-child(5)', '43213423412342134');
  await page.click('#main_body > div > div > form > input[type="submit"]:nth-child(8)');
  
  await browser.close();
})();
