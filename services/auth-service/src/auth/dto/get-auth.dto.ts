import { IsDate, IsString } from "class-validator";

export class GetAuthDto {
    @IsString()
   id: string;
   
   @IsString()
   email: string;
   
   @IsString()
   name?: string;

   @IsDate()
   createdAt: Date;
}