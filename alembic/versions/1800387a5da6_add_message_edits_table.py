"""add message_edits table

Revision ID: 1800387a5da6
Revises: 997c931bcfd6
Create Date: 2025-07-18 13:50:43.126663

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '1800387a5da6'
down_revision: Union[str, None] = '997c931bcfd6'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    pass


def downgrade() -> None:
    """Downgrade schema."""
    pass
