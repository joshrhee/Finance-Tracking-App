import Link from "next/link";
import { DarkModeToggle } from "./DarkModeToggle";

export const Navbar = () => {
    return (
        <div>
            <ul className="flex justify-between my-10 items-center">
                <div>
                    <Link href="/">
                        <li>Home</li>
                    </Link>
                </div>
                <div className="flex gap-10">
                    <DarkModeToggle />
                </div>
            </ul>
        </div>
    );
};
