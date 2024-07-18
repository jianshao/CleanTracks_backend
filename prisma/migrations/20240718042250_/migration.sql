/*
  Warnings:

  - You are about to drop the column `created_time` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `name` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `updated_time` on the `user` table. All the data in the column will be lost.
  - You are about to drop the `members` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `orders` table. If the table is not empty, all the data it contains will be lost.

*/
-- CreateEnum
CREATE TYPE "USER_STATUS" AS ENUM ('INIT', 'SUBSCRIPTED', 'EXPIRED', 'PAUSED', 'CANCELED');

-- AlterTable
ALTER TABLE "user" DROP COLUMN "created_time",
DROP COLUMN "name",
DROP COLUMN "updated_time",
ADD COLUMN     "register_time" TIMESTAMP(3) NOT NULL DEFAULT now(),
ADD COLUMN     "status" "USER_STATUS" NOT NULL DEFAULT 'INIT';

-- DropTable
DROP TABLE "members";

-- DropTable
DROP TABLE "orders";

-- DropEnum
DROP TYPE "MEMBER_STATUS";

-- DropEnum
DROP TYPE "ORDER_STATUS";
