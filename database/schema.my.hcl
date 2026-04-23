schema "tsumiki" {
}

table "users" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "discord_user_id" {
    type = varchar(255)
    null = true
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "guild_id" {
    type = varchar(255)
    null = true
  }
  column "avatar_url" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_discord_user_id" {
    columns = [column.discord_user_id]
    unique  = true
  }
}

table "thumbnails" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "user_id" {
    type = int
    null = false
  }
  column "path" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "fk_thumbnails_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_thumbnails_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "works" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "title" {
    type = varchar(255)
    null = false
  }
  column "description" {
    type = varchar(4095)
    null = true
  }
  column "thumbnail_id" {
    type = int
    null = true
  }
  column "owner_user_id" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "fk_works_owner_user_id" {
    columns = [column.owner_user_id]
  }
  index "fk_works_thumbnail_id" {
    columns = [column.thumbnail_id]
    unique  = true
  }
  foreign_key "fk_works_owner_user_id" {
    columns     = [column.owner_user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_works_thumbnail_id" {
    columns     = [column.thumbnail_id]
    ref_columns = [table.thumbnails.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumikis" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "title" {
    type = varchar(255)
    null = false
  }
  column "thumbnail_id" {
    type = int
    null = true
  }
  column "visibility" {
    type = enum("public", "limited")
    null = false
  }
  column "miruko_watching_duration" {
    type = int
    null = true
  }
  column "work_id" {
    type = int
    null = true
  }
  column "user_id" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "fk_tsumikis_thumbnail_id" {
    columns = [column.thumbnail_id]
    unique  = true
  }
  index "fk_tsumikis_work_id" {
    columns = [column.work_id]
  }
  index "fk_tsumikis_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_tsumikis_thumbnail_id" {
    columns     = [column.thumbnail_id]
    ref_columns = [table.thumbnails.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumikis_work_id" {
    columns     = [column.work_id]
    ref_columns = [table.works.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumikis_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumiki_blocks" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "message" {
    type = varchar(255)
    null = true
  }
  column "percentage" {
    type = int
    null = false
  }
  column "condition" {
    type = int
    null = false
  }
  column "next_block_id" {
    type = int
    null = true
  }
  column "tsumiki_id" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  column "deleted_at" {
    type = timestamp
    null = true
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_tsumiki_blocks_next_block_id" {
    columns = [column.next_block_id]
    unique  = true
  }
  index "idx_tsumiki_blocks_tsumiki_id_next_block_id" {
    columns = [column.tsumiki_id, column.next_block_id]
  }
  foreign_key "fk_tsumiki_blocks_next_block_id" {
    columns     = [column.next_block_id]
    ref_columns = [table.tsumiki_blocks.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_blocks_tsumiki_id" {
    columns     = [column.tsumiki_id]
    ref_columns = [table.tsumikis.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumiki_block_medias" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "type" {
    type = enum("image", "audio", "video")
    null = false
  }
  column "url" {
    type = varchar(255)
    null = false
  }
  column "order" {
    type = int
    null = true
  }
  column "tsumiki_block_id" {
    type = int
    null = true
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_tsumiki_block_medias_tsumiki_block_id_order" {
    columns = [column.tsumiki_block_id, column.order]
    unique  = true
  }
  foreign_key "fk_tsumiki_block_medias_tsumiki_block_id" {
    columns     = [column.tsumiki_block_id]
    ref_columns = [table.tsumiki_blocks.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumiki_favorites" {
  schema = schema.tsumiki
  column "tsumiki_id" {
    type = int
    null = false
  }
  column "user_id" {
    type = int
    null = false
  }
  column "counts" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.tsumiki_id, column.user_id]
  }
  index "fk_tsumiki_favorites_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_tsumiki_favorites_tsumiki_id" {
    columns     = [column.tsumiki_id]
    ref_columns = [table.tsumikis.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_favorites_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumiki_comments" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "tsumiki_id" {
    type = int
    null = false
  }
  column "tsumiki_block_id" {
    type = int
    null = true
  }
  column "user_id" {
    type = int
    null = false
  }
  column "message" {
    type = varchar(255)
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "fk_tsumiki_comments_tsumiki_id" {
    columns = [column.tsumiki_id]
  }
  index "fk_tsumiki_comments_tsumiki_block_id" {
    columns = [column.tsumiki_block_id]
  }
  index "fk_tsumiki_comments_user_id" {
    columns = [column.user_id]
  }
  foreign_key "fk_tsumiki_comments_tsumiki_id" {
    columns     = [column.tsumiki_id]
    ref_columns = [table.tsumikis.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_comments_tsumiki_block_id" {
    columns     = [column.tsumiki_block_id]
    ref_columns = [table.tsumiki_blocks.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_comments_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "reactions" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "shortcode" {
    type = varchar(255)
    null = false
  }
  column "icon_url" {
    type = varchar(255)
    null = false
  }
  column "is_active" {
    type = boolean
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_reactions_shortcode" {
    columns = [column.shortcode]
    unique  = true
  }
}

table "tsumiki_reactions" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "tsumiki_id" {
    type = int
    null = false
  }
  column "user_id" {
    type = int
    null = false
  }
  column "reaction_id" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_tsumiki_reactions_tsumiki_id_user_id_reaction_id" {
    columns = [column.tsumiki_id, column.user_id, column.reaction_id]
    unique  = true
  }
  index "fk_tsumiki_reactions_user_id" {
    columns = [column.user_id]
  }
  index "fk_tsumiki_reactions_reaction_id" {
    columns = [column.reaction_id]
  }
  foreign_key "fk_tsumiki_reactions_tsumiki_id" {
    columns     = [column.tsumiki_id]
    ref_columns = [table.tsumikis.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_reactions_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_reactions_reaction_id" {
    columns     = [column.reaction_id]
    ref_columns = [table.reactions.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "tsumiki_block_reactions" {
  schema = schema.tsumiki
  column "id" {
    type           = int
    null           = false
    auto_increment = true
  }
  column "tsumiki_block_id" {
    type = int
    null = false
  }
  column "user_id" {
    type = int
    null = false
  }
  column "reaction_id" {
    type = int
    null = false
  }
  column "created_at" {
    type    = timestamp
    null    = true
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type      = timestamp
    null      = true
    default   = sql("CURRENT_TIMESTAMP")
    on_update = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_tsumiki_block_reactions_tsumiki_block_id_user_id_reaction_id" {
    columns = [column.tsumiki_block_id, column.user_id, column.reaction_id]
    unique  = true
  }
  index "fk_tsumiki_block_reactions_user_id" {
    columns = [column.user_id]
  }
  index "fk_tsumiki_block_reactions_reaction_id" {
    columns = [column.reaction_id]
  }
  foreign_key "fk_tsumiki_block_reactions_tsumiki_block_id" {
    columns     = [column.tsumiki_block_id]
    ref_columns = [table.tsumiki_blocks.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_block_reactions_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
  foreign_key "fk_tsumiki_block_reactions_reaction_id" {
    columns     = [column.reaction_id]
    ref_columns = [table.reactions.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}
