export interface User {
  ID?: string;
  Name?: string;
  Email?: string;
  Avatar?: string;
}

function UserInfo(user: User) {
  return (
    <div className="flex flex-row gap-5 items-center">
      <img src={user.Avatar} className="rounded-full w-12 h-12" />
      <div className="flex flex-col">
        <p className="text-2xl font-semibold">{user.Name}</p>
        <p className="text-lg text-gray-800">{user.Email}</p>
      </div>
    </div>
  );
}

export default UserInfo;
