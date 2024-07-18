/*
  Warnings:

  - You are about to drop the column `status` on the `user` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "SUB_STATUS" AS ENUM ('INIT', 'PAID', 'EXPIRED', 'PAUSED', 'CANCELLED');

-- AlterTable
ALTER TABLE "user" DROP COLUMN "status",
ALTER COLUMN "register_time" SET DEFAULT now();

-- DropEnum
DROP TYPE "USER_STATUS";

-- CreateTable
CREATE TABLE "subscriptions" (
    "id" SERIAL NOT NULL,
    "uid" INTEGER NOT NULL,
    "store_id" INTEGER NOT NULL,
    "product_id" INTEGER NOT NULL,
    "variant_id" INTEGER NOT NULL,
    "subscription_id" INTEGER NOT NULL,
    "status" "SUB_STATUS" NOT NULL DEFAULT 'INIT',
    "created_time" TIMESTAMP(3) NOT NULL DEFAULT now(),

    CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id")
);
