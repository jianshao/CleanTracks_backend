-- CreateEnum
CREATE TYPE "MEMBER_STATUS" AS ENUM ('USING', 'DISUSE');

-- CreateEnum
CREATE TYPE "ORDER_STATUS" AS ENUM ('CREATED', 'PAID');

-- CreateTable
CREATE TABLE "user" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL DEFAULT '',
    "avatar" TEXT NOT NULL DEFAULT '',
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "created_time" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_time" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "members" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "desc" TEXT NOT NULL DEFAULT '',
    "type" INTEGER NOT NULL DEFAULT 0,
    "status" "MEMBER_STATUS" NOT NULL DEFAULT 'USING',
    "created_time" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "updated_time" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "members_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "orders" (
    "id" SERIAL NOT NULL,
    "uid" INTEGER NOT NULL,
    "member_id" INTEGER NOT NULL,
    "status" "ORDER_STATUS" NOT NULL DEFAULT 'CREATED',
    "created_time" TIMESTAMP(3) NOT NULL DEFAULT now(),
    "paid_time" TIMESTAMP(3) NOT NULL,
    "expire_time" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "orders_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON "user"("email");

-- CreateIndex
CREATE UNIQUE INDEX "members_name_key" ON "members"("name");
