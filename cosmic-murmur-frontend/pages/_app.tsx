import '../styles/globals.css'
import type { NextPageWithLayout } from 'types'
import type { AppProps } from 'next/app'

import NavLayout from 'layouts/nav-layout'

import Head from 'next/head'

type AppPropsWithLayout = AppProps & {
  Component: NextPageWithLayout
}

function App({ Component, pageProps }: AppPropsWithLayout) {

  const getLayout = Component.getLayout
      ? Component.getLayout
      : NavLayout

  return (
      <>
        <Head>
          <title>Cosmic Murmur</title>
        </Head>
        { getLayout(<Component {...pageProps} />) }
      </>
  )
}

export default App
