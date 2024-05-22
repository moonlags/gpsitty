import { Locate, User } from "lucide-react";
import UserInfo, { IUser } from "./UserInfo";
import { Button } from "./ui/button";

const handleLogin = () => {
  window.location.href =
    "http://" + import.meta.env.VITE_BACKEND_HOST + "/auth/google";
};

function NavBar(user: IUser) {
  return (
    <div className="px-10 items-center flex flex-row justify-between w-full h-16 bg-gray-300">
      <div className="flex flex-row gap-2 items-center">
        <Locate className="w-8 h-8" />
        <p className="text-xl">GPSitty</p>
      </div>
      <div className="flex flex-row gap-10">
        {user.ID ? (
          <UserInfo {...user} />
        ) : (
          <Button className="flex flex-row gap-2 dark" onClick={handleLogin}>
            <User />
            Login With Google
          </Button>
        )}
      </div>
    </div>
  );
}

export default NavBar;
