
import type {ReactElement} from "react";
import type {NextPageWithLayout} from "types";

import CenteredLayout from 'layouts/centered-layout'


const Home: NextPageWithLayout = () => {
  return (
    <div>
      THIS IS A LOGIN PAGE
    </div>
  )
}

Home.getLayout = (page: ReactElement) => {
    return CenteredLayout(page)
}

export default Home
