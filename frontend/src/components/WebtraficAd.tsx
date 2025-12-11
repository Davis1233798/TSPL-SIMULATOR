import React, { useEffect } from 'react';

const WebtraficAd: React.FC = () => {
    useEffect(() => {
        // Check if script already exists to prevent duplicates
        if (document.querySelector('script[src="https://webtrafic.ru/ads.php?uid=17708"]')) {
            return;
        }

        const script = document.createElement('script');
        script.src = "https://webtrafic.ru/ads.php?uid=17708";
        script.async = true;

        const container = document.getElementById('webtraf_17708');
        if (container) {
            container.appendChild(script);
        }

        return () => {
            // Cleanup if needed, though ad scripts usually shouldn't be removed aggressively
            // keeping it simple for now
        };
    }, []);

    return (
        <div
            id="webtraf_17708"
            style={{
                width: '468px',
                height: '60px',
                // Add some margin if needed for mobile/responsiveness
            }}
            className="webtrafic-ad"
        />
    );
};

export default WebtraficAd;
