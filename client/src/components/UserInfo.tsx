export interface IUser {
  ID?: string;
  Name?: string;
  Email?: string;
  Avatar?: string;
}

function UserInfo(user: IUser) {
  return (
    <div className="flex flex-row gap-5 items-center">
      <img
        src={user.Avatar}
        className="rounded-full w-10 h-10 border border-gray-700"
      />
      <div className="flex flex-col">
        <p className="text-xl">{user.Name}</p>
        <p className="text-md text-gray-800">{user.Email}</p>
      </div>
    </div>
  );
}

export default UserInfo;
