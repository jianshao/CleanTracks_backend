/*
  Warnings:

  - A unique constraint covering the columns `[subscription_id]` on the table `subscriptions` will be added. If there are existing duplicate values, this will fail.

*/
-- AlterTable
ALTER TABLE "subscriptions" ALTER COLUMN "created_time" SET DEFAULT now();

-- AlterTable
ALTER TABLE "user" ALTER COLUMN "register_time" SET DEFAULT now();

-- CreateIndex
CREATE UNIQUE INDEX "subscriptions_subscription_id_key" ON "subscriptions"("subscription_id");
