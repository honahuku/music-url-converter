// downloadPage/main.js
const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({
    headless: "new",
    args: ['--lang=ja']
  });
  const page = await browser.newPage();
  page.setUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36");
  process.stdout.write("waiting network...");
  await page.goto(process.argv[2], { waitUntil: 'networkidle0' }); // ネットワークがアイドル状態になるまで待機
  process.stdout.write("Done!\n");
  await page.setExtraHTTPHeaders({
    'Accept-Language': 'ja-JP'
  });

  // XPathが表示されるのを待機
  process.stdout.write("waiting for XPath...");
  const xpathExpression = '/html/body/ytmusic-app/ytmusic-app-layout';
  await page.waitForXPath(xpathExpression, { visible: true });
  process.stdout.write("Done!\n");

  // トラック名を取得
  const trackNameXPath = '/html/body/ytmusic-app/ytmusic-app-layout/ytmusic-player-bar/div[2]/div[2]/yt-formatted-string';
  await page.waitForXPath(trackNameXPath, { visible: true });
  const [trackElement] = await page.$x(trackNameXPath);
  const trackName = await page.evaluate((el: Element) => el.getAttribute('title'), trackElement);

  // アーティスト名を取得
  const artistNameXPath = '/html/body/ytmusic-app/ytmusic-app-layout/ytmusic-player-bar/div[2]/div[2]/span/span[2]/yt-formatted-string/a[1]';
  await page.waitForXPath(artistNameXPath, { visible: true });
  const [artistElement] = await page.$x(artistNameXPath);
  const artistName = await page.evaluate((el: Element) => el.textContent, artistElement);

  // トラック名とアーティスト名を標準出力に出力
  console.log(JSON.stringify({ trackName, artistName }));

  await browser.close();
})();
