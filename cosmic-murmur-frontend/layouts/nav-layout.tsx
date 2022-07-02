
import type {ReactElement} from "react";

export default function Layout(element: ReactElement) {
    return (
        <section>
            <header>
                nav
            </header>
            <main>
                { element }
            </main>
        </section>

    )
}
