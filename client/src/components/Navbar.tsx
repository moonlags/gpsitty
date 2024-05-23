import UserInfo, { IUser } from "./UserInfo";
import LoginButton from "./LoginButton";
import { IoLocate } from "solid-icons/io";

function NavBar(props: { user: IUser | undefined }) {
  return (
    <div class="px-10 items-center flex flex-row justify-between w-full min-h-16 bg-gray-300">
      <div class="flex flex-row gap-2 items-center">
        <IoLocate class="w-8 h-8" />
        <p class="text-xl">GPSitty</p>
      </div>
      <div class="flex flex-row gap-10">
        {props.user ? <UserInfo user={props.user} /> : <LoginButton />}
      </div>
    </div>
  );
}

export default NavBar;
