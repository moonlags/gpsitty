import { AiOutlineUser } from "solid-icons/ai";
import { ImCog } from "solid-icons/im";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import axios from "axios";

export interface IUser {
  ID?: string;
  Email?: string;
}

function UserInfo(props: { user: IUser }) {
  function handleLogout() {
    axios
      .get(import.meta.env.VITE_BACKEND_HOST + "/auth/logout")
      .then(() => {
        location.reload();
      })
      .catch((err) => {
        console.error(err);
      });
  }

  return (
    <div class="flex flex-row gap-8 items-center">
      <div class="flex flex-row gap-1">
        <AiOutlineUser class="h-6 w-6" />
        <p class="text-md text-gray-800">{props.user.Email}</p>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <ImCog />
        </DropdownMenuTrigger>
        <DropdownMenuContent class="w-48">
          <DropdownMenuItem onClick={handleLogout}>
            <span>Logout</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}

export default UserInfo;
