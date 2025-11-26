const fs = require('fs');
const path = require('path');
const https = require('https');
const crypto = require('crypto');

// Configuration
const AD_BLOCK_URL = 'https://adbpage.com/adblock?v=3&format=js';
const PUBLIC_DIR = path.join(__dirname, '../public');
const INDEX_HTML_PATH = path.join(PUBLIC_DIR, 'index.html');

// Generate a random filename (e.g., ab_x7z8.js)
function generateRandomFilename() {
    const randomString = crypto.randomBytes(4).toString('hex');
    return `ab_${randomString}.js`;
}

// Download the file
function downloadFile(url, dest) {
    return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(dest);
        https.get(url, (response) => {
            if (response.statusCode !== 200) {
                reject(new Error(`Failed to download: ${response.statusCode}`));
                return;
            }
            response.pipe(file);
            file.on('finish', () => {
                file.close();
                resolve();
            });
        }).on('error', (err) => {
            fs.unlink(dest, () => { }); // Delete the file async
            reject(err);
        });
    });
}

// Update index.html
function updateIndexHtml(newScriptFilename) {
    let htmlContent = fs.readFileSync(INDEX_HTML_PATH, 'utf8');

    // Regex to find existing adblock script tags (assuming they follow a pattern or we track them)
    // Since we want to replace the one we added, we can look for the specific comment or structure
    // But for robustness, let's look for the script tag with the specific source pattern if possible,
    // or just manage a specific marker.

    // Strategy: Look for a marker comment. If not found, insert before </head>.
    // If found, replace the script tag following it.

    const MARKER_START = '<!-- ANTI-ADBLOCK-START -->';
    const MARKER_END = '<!-- ANTI-ADBLOCK-END -->';
    const scriptTag = `<script id="aclib" type="text/javascript" src="%PUBLIC_URL%/${newScriptFilename}"></script>`;
    const newBlock = `${MARKER_START}\n    ${scriptTag}\n    ${MARKER_END}`;

    const regex = new RegExp(`${MARKER_START}[\\s\\S]*?${MARKER_END}`, 'g');

    if (regex.test(htmlContent)) {
        console.log('Updating existing Anti-Adblock script tag...');
        htmlContent = htmlContent.replace(regex, newBlock);
    } else {
        console.log('Inserting new Anti-Adblock script tag...');
        // Insert before </head>
        htmlContent = htmlContent.replace('</head>', `${newBlock}\n  </head>`);
    }

    fs.writeFileSync(INDEX_HTML_PATH, htmlContent);
}

// Clean up old files
function cleanOldFiles(currentFilename) {
    fs.readdir(PUBLIC_DIR, (err, files) => {
        if (err) {
            console.error('Error reading public directory:', err);
            return;
        }

        files.forEach(file => {
            if (file.startsWith('ab_') && file.endsWith('.js') && file !== currentFilename) {
                const filePath = path.join(PUBLIC_DIR, file);
                fs.unlink(filePath, (err) => {
                    if (err) console.error(`Failed to delete old file ${file}:`, err);
                    else console.log(`Deleted old file: ${file}`);
                });
            }
        });
    });
}

async function main() {
    try {
        console.log('Starting Anti-Adblock update...');

        const newFilename = generateRandomFilename();
        const destPath = path.join(PUBLIC_DIR, newFilename);

        console.log(`Downloading library to ${newFilename}...`);
        await downloadFile(AD_BLOCK_URL, destPath);

        console.log('Updating index.html...');
        updateIndexHtml(newFilename);

        console.log('Cleaning up old files...');
        cleanOldFiles(newFilename);

        console.log('Anti-Adblock update complete!');
    } catch (error) {
        console.error('Error updating Anti-Adblock:', error);
        process.exit(1);
    }
}

main();
