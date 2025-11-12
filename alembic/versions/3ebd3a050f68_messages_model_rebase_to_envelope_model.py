"""messages model rebase to envelope model

Revision ID: 3ebd3a050f68
Revises: ea3a5800ab6f
Create Date: 2025-11-10 18:02:25.478595
"""

from typing import Sequence, Union
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '3ebd3a050f68'
down_revision: Union[str, None] = 'ea3a5800ab6f'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.drop_table('message_edits')
    op.drop_table('attachments')
    op.drop_table('messages')

    op.create_table(
        'messages',
        sa.Column('id', sa.String(), primary_key=True),
        sa.Column('chat_id', sa.Integer(), sa.ForeignKey('chats.id', ondelete='CASCADE'), nullable=False, index=True),
        sa.Column('sender_id', sa.Integer(), sa.ForeignKey('users.id', ondelete='SET NULL'), nullable=True, index=True),
        sa.Column('target_id', sa.Integer(), sa.ForeignKey('users.id', ondelete='SET NULL'), nullable=True, index=True),

        sa.Column('kind', sa.String(), nullable=False, server_default='message'),
        sa.Column('action', sa.String(), nullable=False, server_default='send'),
        sa.Column('chat_type', sa.String(), nullable=False, server_default='private'),
        sa.Column('timestamp', sa.BigInteger(), nullable=True),

        sa.Column('text', sa.Text(), nullable=True),
        sa.Column('attachments', sa.JSON(), nullable=True),
        sa.Column('reply_to_id', sa.String(), sa.ForeignKey('messages.id', ondelete='SET NULL'), nullable=True),

        sa.Column('is_edited', sa.Boolean(), server_default=sa.text('false')),
        sa.Column('deleted', sa.Boolean(), server_default=sa.text('false')),
        sa.Column('is_pinned', sa.Boolean(), server_default=sa.text('false')),

        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.func.now()),
        sa.Column('edited_at', sa.DateTime(timezone=True), nullable=True),
        sa.Column('raw_envelope', sa.JSON(), nullable=True),
    )

    op.create_index('idx_messages_chat_created', 'messages', ['chat_id', 'created_at'])

    op.create_table(
        'attachments',
        sa.Column('id', sa.Integer(), primary_key=True, autoincrement=True),
        sa.Column('message_id', sa.String(), sa.ForeignKey('messages.id', ondelete='CASCADE'), nullable=False, index=True),
        sa.Column('filename', sa.String(), nullable=False),
        sa.Column('filetype', sa.String(), nullable=False),
        sa.Column('filesize', sa.Integer(), nullable=False),
        sa.Column('original_name', sa.String(), nullable=True),
        sa.Column('uploaded_at', sa.DateTime(timezone=True), server_default=sa.func.now())
    )

    op.create_table(
        'message_edits',
        sa.Column('id', sa.Integer(), primary_key=True, autoincrement=True),
        sa.Column('message_id', sa.String(), sa.ForeignKey('messages.id', ondelete='CASCADE'), nullable=False, index=True),
        sa.Column('old_text', sa.Text(), nullable=False),
        sa.Column('edited_at', sa.DateTime(timezone=True), server_default=sa.func.now())
    )


def downgrade() -> None:
    op.drop_table('message_edits')
    op.drop_table('attachments')
    op.drop_table('messages')