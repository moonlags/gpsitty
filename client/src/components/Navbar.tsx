import UserInfo, { User } from "./UserInfo";

const handleLogin = () => {
  window.location.href = "http://localhost:50731/auth/google";
};

function NavBar(user: User) {
  return (
    <div className="px-10 items-center flex flex-row justify-between w-full h-16 bg-gray-300">
      <p className="text-3xl font-semibold">GPSitty</p>
      <div className="flex flex-row gap-10">
        {user.ID ? (
          <UserInfo {...user} />
        ) : (
          <button
            className="bg-gray-200 px-10 py-3 rounded-md shadow-md font-bold hover:scale-110 duration-100"
            onClick={handleLogin}
          >
            Login
          </button>
        )}
      </div>
    </div>
  );
}

export default NavBar;
