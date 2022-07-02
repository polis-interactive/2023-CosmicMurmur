
import {ReactElement, ReactNode} from "react";
import {NextPage} from "next";

export type WithLayout = { getLayout?: (page: ReactElement) => ReactNode }

export type NextPageWithLayout = NextPage & WithLayout
