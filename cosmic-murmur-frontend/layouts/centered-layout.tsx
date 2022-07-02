
import type {ReactElement} from "react";

export default function Layout(element: ReactElement) {
    return (
        <section>
            <header>
                centered
            </header>
            <main>
                { element }
            </main>
        </section>

    )
}
