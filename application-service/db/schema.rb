# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.1].define(version: 0) do
  create_table "applications", id: { type: :bigint, unsigned: true }, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.string "name", null: false
    t.string "token", null: false
    t.integer "chats_count", default: 0, null: false
    t.timestamp "created_at", null: false
    t.timestamp "updated_at", null: false
    t.index ["token"], name: "unique_token", unique: true
  end

  create_table "chats", id: { type: :bigint, unsigned: true }, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.string "application_id", null: false  # application_id references token, not id
    t.integer "number", null: false
    t.integer "messages_count", default: 0, null: false
    t.timestamp "created_at", null: false
    t.index ["application_id", "number"], name: "unique_app_number", unique: true
  end

  create_table "messages", id: { type: :bigint, unsigned: true }, charset: "utf8mb4", collation: "utf8mb4_0900_ai_ci", force: :cascade do |t|
    t.bigint "chat_id", null: false, unsigned: true
    t.integer "number", null: false
    t.text "body", null: false
    t.timestamp "created_at", null: false
    t.index ["chat_id", "number"], name: "unique_chat_number", unique: true
  end

  # Add foreign key constraint referencing token in applications
  add_foreign_key "chats", "applications", column: "application_id", primary_key: "token", name: "chats_ibfk_1"

  add_foreign_key "messages", "chats", column: "chat_id", primary_key: "id", name: "messages_ibfk_1"
end
