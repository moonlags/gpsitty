import deviceImg from "../assets/device.jpg";
import { BsBatteryCharging, BsBatteryFull } from "solid-icons/bs";
import { Button } from "./ui/button";

export interface IDevice {
  Imei: string;
  BatteryPower: number;
  Charging: boolean;
}

function DeviceCard(props: { device: IDevice }) {
  return (
    <div class="flex flex-col gap-6 border rounded-lg shadow-lg p-8 justify-around items-center">
      <img src={deviceImg} />
      <div class="flex flex-row justify-between w-full items-center">
        <p>IMEI: {props.device.Imei}</p>
        <div class="flex flex-row gap-1 items-center">
          <p>{props.device.BatteryPower}%</p>
          {props.device.Charging ? <BsBatteryCharging /> : <BsBatteryFull />}
        </div>
      </div>
      <Button class="dark">Open Map</Button>
    </div>
  );
}

export default DeviceCard;
