import {
  pgTable,
  uuid,
  text,
  timestamp,
  integer,
  boolean,
} from "drizzle-orm/pg-core";
import { relations, sql } from "drizzle-orm";

export const files = pgTable("files", {
  id: uuid("id").primaryKey().defaultRandom(),

  // basic file/folder info
  name: text("name").notNull(),
  path: text("path").notNull(), // path to the file like /home/user/file.txt
  size: integer("size").notNull(),
  type: text("type").notNull(), // MIME type of the file or folder

  // storage info
  fileUrl: text("file_url").notNull(), // URL to the file or folder
  thumbnailUrl: text("thumbnail_url"),

  // ownership info
  userId: text("user_id").notNull(), // Clerk user ID
  parentId: uuid("parent_id"), // ID of the parent folder (null if root folder)

  // file/folder flags
  isFolder: boolean("is_folder").notNull().default(false),
  isStarred: boolean("is_starred").notNull().default(false),
  isTrash: boolean("is_trash").notNull().default(false),
  isShared: boolean("is_shared").notNull().default(false),

  // file/folder metadata
  createdAt: timestamp("created_at").notNull().defaultNow(),
  updatedAt: timestamp("updated_at")
    .notNull()
    .default(sql`CURRENT_TIMESTAMP`)
    .$onUpdate(() => sql`(CURRENT_TIMESTAMP)`),
});

export const filesRelations = relations(files, ({ one, many }) => ({
  // parent: Each file/folder can have one parent folder
  // relationship to parent folder
  parent: one(files, {
    fields: [files.parentId],
    references: [files.id],
  }),

  // children: Each file/folder can have many child files/folders
  // relationship to child file/folder
  children: many(files),
}));

export const File = typeof files.$inferSelect;
export const NewFile = typeof files.$inferInsert;
