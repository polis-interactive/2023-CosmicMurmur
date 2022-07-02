import { Html, Head, Main, NextScript } from 'next/document'

export default function Document() {
    return (
        <Html>
            <Head>
                <link rel="preconnect" href="https://fonts.googleapis.com" />
                <link
                    href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,100;0,400;0,700;1,100;1,400;1,700&display=swap"
                    rel="stylesheet"
                />
                <meta name="description" content="Cosmic Murmur, an installation by sparky, coded with love by broose" />
                <link rel="icon" href="/images/site/favicon.ico" />
            </Head>
            <body>
            <Main />
            <NextScript />
            </body>
        </Html>
    )
}
