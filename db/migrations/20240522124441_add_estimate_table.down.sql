ALTER TABLE "estimateOrder" DROP CONSTRAINT IF EXISTS "estimateOrder_estimateId_fkey";
ALTER TABLE "estimateOrder" DROP CONSTRAINT IF EXISTS "estimateOrder_merchantId_fkey";
ALTER TABLE "estimateOrderItem" DROP CONSTRAINT IF EXISTS "estimateOrderItem_estimateOrderId_fkey";
ALTER TABLE "estimateOrderItem" DROP CONSTRAINT IF EXISTS "estimateOrderItem_itemId_fkey";

DROP TABLE IF EXISTS "estimate";
DROP TABLE IF EXISTS "estimateOrder";
DROP TABLE IF EXISTS "estimateOrderItem";